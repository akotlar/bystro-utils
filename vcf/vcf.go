package vcf

import (
  "strconv"
  "strings"
  "bytes"
)

const ChromIdx int = 0
const PosIdx int = 1
const IdIdx int = 2
const RefIdx int = 3
const AltIdx int = 4
const QualIdx int = 5
const FilterIdx int = 6
const InfoIdx int = 7
const FormatIdx int = 8

func UpdateFieldsWithAlt(ref string, alt string, pos string, multiallelic bool) (string, string, string, string, error) {
  /*********************** SNPs *********************/
  if len(alt) == len(ref) {
    if alt == ref {
      // No point in returning ref sites
      return "", "", "", "", nil
    }

    if len(ref) > 1 {
      // SNPs that are multiallelic with indels, can be longer than 1 base long
      // Confusingly enough, so can MNPs
      // So, we will check for both
      if ref[0] != alt[0] && ref[1] != alt[1]{
        // This is most likely an MNP
        // As MNPs are contiguous
        // Currently we haven't enabled MNP parsing in the caller
        return "", "", "", "", nil
      }

      diffIdx := -1
      // Let's check each base; if there is more than 1 change, that is an error
      for i := 0; i < len(ref); i++ {
        if ref[i] != alt[i] {
          // Major red flag, there should never be a len(ref) == len(alt) allele that isn't an MNP or SNP
          // TODO: should we relax this? If we allow MNP, may as well allow sparse MNPs
          if diffIdx > -1 {
            return "", "", "", "", nil
          }

          diffIdx = i
        }
      }

      // Most cases are diffIdx == 0, allow us to skip 1 strconv.Atoi, 1 assignment, 1 strconv.Itoa, and 1 addition
      if diffIdx == 0 {
        return "SNP", pos, string(ref[diffIdx]), string(alt[diffIdx]), nil
      }

      intPos, _ := strconv.Atoi(pos)
      return "SNP", strconv.Itoa(intPos + diffIdx), string(ref[diffIdx]), string(alt[diffIdx]), nil
    }

    return "SNP", pos, ref, alt, nil
  }

  /*********************** INSERTIONS AND DELETIONS *********************/
  // TODO: Handle case where first base of contig is deleted, and padded as the first unmodified base downstream
  //First base is always padding
  if ref[0] != alt[0] {
    return "", "", "", "", nil
  }

  /*************************** DELETIONS FIRST **************************/
  if len(ref) > len(alt) {
    intPos, err := strconv.Atoi(pos)
    if err != nil {
      return "", "", "", "", err
    }

    // Simple insertions delete the entire reference, sans the padding base to the left
    // TODO: handle 1st base deleted in contig, padded to right
    if len(alt) == 1 {
      return "DEL", strconv.Itoa(intPos + 1), string(ref[1]), strconv.Itoa(1 - len(ref)), nil
    }

    // log.Println("Complex del", pos, ref, alt, multiallelic)

    // Complex deletions, inside of a reference
    // Ex: Ref: TCT Alt: T, TT (the TT is a 1 base C deletion)
    //this typically only comes up with multiallelics that have a 2nd deletion, that covers all bases (excepting 1 padding base) in the reference
    //In other cases, just skip for now, mostly seems like an error
    if multiallelic == false {
      return "", "", "", "", nil
    }

    // Our deletion should happen within the reference, so the non-padded
    // portion of the reference is what we'll check
    if strings.Contains(ref, alt[1: ]) == false {
      return "", "", "", "", nil
    }
    // TODO: More precise checking; for instance we can check if the alt is contained within the end of the ref (sans the 1 base deletion)
    return "DEL", strconv.Itoa(intPos + 1), string(ref[1]), strconv.Itoa(len(alt) - len(ref)), nil
  }

  /*********************** INSERTIONS *********************/
  // len(ref) > 1 should always indicate that this is a multiallelic that contains a deletion
  // therefore requiring 1 base of padding to the left
  // there may be cases where VCF variants are unnecessarily padded, but lets skip these
  if len(ref) > 1 {
    if multiallelic == false {
      return "", "", "", "", nil
    }

    // log.Println("Complex ins", pos, ref, alt, multiallelic)

    // Our insertion should happen within the reference, so the non-padded
    // portion of the reference is what we'll check
    if strings.Contains(alt, ref[1: ]) == false {
      return "", "", "", "", nil
    }
    // TODO: More precise checking; for instance we can check if the alt is contained within the end of the ref (sans the 1 base deletion)
    var insBuffer bytes.Buffer
    insBuffer.WriteString("+")
    insBuffer.WriteString(alt[1:len(alt) - len(ref) + 1])
    return "INS", pos, string(ref[0]), insBuffer.String(), nil
  }

  var insBuffer bytes.Buffer
  insBuffer.WriteString("+")
  insBuffer.WriteString(alt[1:])
  return "INS", pos, ref, insBuffer.String(), nil
}