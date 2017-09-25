package parse

import (
  "bufio"
  "regexp"
)

const (
  NotTrTv string = "0"
  Tr string = "1"
  Tv string = "2"
)

const (
  Snp string = "SNP"
  Ins string = "INS"
  Del string = "DEL"
  Mnp string = "MNP"
  Multi string = "MULTIALLELIC"
  Dsnp string = "DENOVO_SNP"
  Dins string = "DENOVO_INS"
  Ddel string = "DENOVO_DEL"
  Dmulti string = "DENOVO_MULTIALLELIC"
)

var Header = []string{"chrom", "pos", "type", "ref", "alt", "trTv", "heterozygotes",
    "heterozygosity", "homozygotes", "homozygosity", "missingGenos", "missingness", "sampleMaf"}

// Note: If passed decomposed multiallelics, will happily consider them tr or tv
// I believe it should be up to the consumer what to do with this indeterminate case
func GetTrTv(ref string, alt string) string {
  if ref == "A" {
    if alt == "G" {
      return Tr
    }

    if alt == "C" || alt == "T" {
      return Tv
    }

    return NotTrTv
  }

  if ref == "G" {
    if alt == "A" {
      return Tr
    }

    if alt == "C" || alt == "T" {
      return Tv
    }

    return NotTrTv
  }

  if ref == "C" {
    if alt == "T" {
      return Tr
    }

    if alt == "A" || alt == "G" {
      return Tv
    }

    return NotTrTv
  }

  if ref == "T" {
    if alt == "C" {
      return Tr
    }

    if alt == "A" || alt == "G" {
      return Tv
    }

    return NotTrTv
  }

  return NotTrTv
}

func NormalizeHeader(header []string) {
  re := regexp.MustCompile(`[^a-zA-Z0-9\_\-\#]`)

  for i := 0; i < len(header); i+= 1 {
    if len(header[i]) > 0 {
      header[i] = re.ReplaceAllString(header[i], "_")
    }
  }
}

func FindEndOfLine (r *bufio.Reader, s string) (byte, int, string, error) {
  runeChar, _, err := r.ReadRune()

  if err != nil {
    return byte(0), 0, "", err
  }

  if runeChar == '\r' {
    nextByte, err := r.Peek(1)

    if err != nil {
      return byte(0), 0, "", err
    }

    if rune(nextByte[0]) == '\n' {
      //Remove the line feed
      _, _, err = r.ReadRune()

      if err != nil {
        return byte(0), 0, "", err
      }

      return nextByte[0], 2, s, nil
    }

    return byte('\r'), 1, s, nil
  }

  if runeChar == '\n' {
    return byte('\n'), 1, s, nil
  }

  s += string(runeChar)
  return FindEndOfLine(r, s)
}

func AppendMissing(numAlt int, sampleName string, arr [][]string) {
  for i := 0; i < numAlt; i++ {
    arr[i] = append(arr[i], sampleName)
  }
}