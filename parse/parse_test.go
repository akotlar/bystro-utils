package parse

import (
  "testing"
  "strings"
)

func TestTrTv(t *testing.T) {
  if GetTrTv("A", "T", false) != '2' || GetTrTv("T", "A", false) != '2' ||
  GetTrTv("A", "C", false) != '2' || GetTrTv("C", "A", false) != '2' ||
  GetTrTv("G", "C", false) != '2' || GetTrTv("C", "G", false) != '2' ||
  GetTrTv("G", "T", false) != '2' || GetTrTv("T", "G", false) != '2' {
    t.Error("Couldn't parse transversions")
  }

  if GetTrTv("A", "G", false) != '1' || GetTrTv("G", "A", false) != '1' ||
  GetTrTv("C", "T", false) != '1' || GetTrTv("T", "C", false) != '1' {
    t.Error("Couldn't parse transversions")
  }

  if GetTrTv("A", "-1", false) != '0' || GetTrTv("A", "+A", false) != '0' {
    t.Error("Couldn't parse non-TrTv sites")
  }

  if GetTrTv("A", "T,C", false) != '0' || GetTrTv("A", "T", true) != '0' {
    t.Error("Couldn't parse non-TrTv sites due to multiallelic")
  }

  if GetTrTv("A", "A", false) != '0' {
    t.Error("Couldn't parse non-TrTv sites that are mistakenly homozygous reference")
  }
}

func TestNormalizationOfSamples(t *testing.T) {
  header := []string{"#CHROM", "POS", "ID", "REF", "ALT", "QUAL", "FILTER", "INFO", "FORMAT", "S1.HAHAHAH", "S2.TRYINGTO.MESSYOUUP", "S3", "S-4"}

  NormalizeHeader(header)

  for i := 9; i < len(header); i++ {
    if strings.Contains(header[i], ".") {
      t.Error("NOT OK: Couldn't replace period")
    } else {
      t.Log("OK: no periods found in", header[i])
    }
  }

  if header[9] == "S1_HAHAHAH" {
    t.Log("OK: replaced period in S1.HAHAHAH", header[9])
  } else {
    t.Error("NOT OK: Couldn't replace period in S1.HAHAHAH", header[9])
  }

  if header[10] == "S2_TRYINGTO_MESSYOUUP" {
    t.Log("OK: replaced two periods in S2.TRYINGTO.MESSYOUUP", header[10])
  } else {
    t.Error("NOT OK: Couldn't replace periods in S2.TRYINGTO.MESSYOUUP", header[10])
  }

  if header[11] == "S3" {
    t.Log("OK: didn't mess up name S3", header[11])
  } else {
    t.Error("NOT OK: Messed up name S3", header[11])
  }

  if header[12] == "S-4" {
    t.Log("OK:  didn't mess up name without a period", header[12])
  } else {
    t.Error("NOT OK: Messed up name S-4", header[12])
  }
}