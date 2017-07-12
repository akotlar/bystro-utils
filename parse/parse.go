package parse

import (
  "bufio"
  "regexp"
)

func TrTv(ref string, alt string) rune {
  if len(alt) > 1 {
    return '0'
  }

  // Transition
  if (ref == "A" && alt == "G") || (ref == "G" && alt == "A") || (ref == "C" && alt =="T") ||
  (ref == "T" && alt == "C") {
    return '1'
  }

  // Transversion
  return '2'
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