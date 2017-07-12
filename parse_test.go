package parse

func TestTrTv(t *testing.T) {
  if trTvStatus("A", "T") != '2' || trTvStatus("T", "A") != '2' ||
  trTvStatus("A", "C") != '2' || trTvStatus("C", "A") != '2' ||
  trTvStatus("G", "C") != '2' || trTvStatus("C", "G") != '2' ||
  trTvStatus("G", "T") != '2' || trTvStatus("T", "G") != '2' {
    t.Error("Couldn't parse transversions")
  }

  if trTvStatus("A", "G") != '1' || trTvStatus("G", "A") != '1' ||
  trTvStatus("C", "T") != '1' || trTvStatus("T", "C") != '1' {
    t.Error("Couldn't parse transversions")
  }

  if trTvStatus("A", "-1") != '0' || trTvStatus("A", "+A") != '0' {
    t.Error("Couldn't parse non-trTv sites")
  }
}