package main

import (
  "fmt"
  "os"
  "encoding/csv"
  "io"
  "log"
  "strconv"
  "time"

  dataframe "github.com/go-gota/gota/dataframe"
  series "github.com/go-gota/gota/series"
  civil "cloud.google.com/go/civil"
)

func parseIntDefaultZero(s string) int {
  res, err := strconv.Atoi(s)
  if err != nil {
    return 0
  } 
  return res
}

func parseDateNum(s string) civil.Date {
  year,err := strconv.Atoi(s[:4])
  if err != nil { panic(err) }
  month,err := strconv.Atoi(s[4:6])
  if err != nil { panic(err) }
  day,err := strconv.Atoi(s[6:8])
  if err != nil { panic(err) }
  return civil.Date{year, time.Month(month), day}
}

func readIncidencePoland() dataframe.DataFrame {
  f, err := os.Open("./data/data-pl.csv")
  if err != nil {
    log.Fatal(err)
  }

  reader := csv.NewReader(f)
  reader.Comma = ';'
  reader.FieldsPerRecord = -1

  dateNum := series.New([]string{}, series.String, "DateNum")
  county := series.New([]string{}, series.String, "County")
  city := series.New([]string{}, series.String, "City")
  cases := series.New([]int{}, series.Int, "Cases")

  for {
        record, err := reader.Read()

        if err == io.EOF {
            break
        }

        if err != nil {
            log.Fatal(err)
        }

        dateNum.Append(record[0])
        county.Append(record[1])
        city.Append(record[2])
        cases.Append(parseIntDefaultZero(record[3]))
  }

  return dataframe.New(dateNum, county, city, cases)
}

func readIncidenceGermany() dataframe.DataFrame {
  f, err := os.Open("./data/data-de.csv")
  if err != nil {
    log.Fatal(err)
  }

  reader := csv.NewReader(f)
  reader.FieldsPerRecord = -1

  date := series.New([]string{}, series.String, "Date")
  county := series.New([]string{}, series.String, "County")
  city := series.New([]string{}, series.String, "City")
  cases := series.New([]int{}, series.Int, "Cases")
  caseFreshness := series.New([]int{}, series.Int, "CaseFreshness")

  // skip header
  reader.Read()

  for {
        record, err := reader.Read()

        if err == io.EOF {
            break
        }

        if err != nil {
            log.Fatal(err)
        }

        date.Append(record[8])
        county.Append(record[2])
        city.Append(record[3])
        cases.Append(parseIntDefaultZero(record[6]))
        caseFreshness.Append(parseIntDefaultZero(record[11]))
  }

  return dataframe.New(date, county, city, cases, caseFreshness)
}

func main() {
  poland := readIncidencePoland()
  germany := readIncidenceGermany()

  berlin := germany.Filter(
    dataframe.F{
      Colname:    "County",
      Comparator: series.Eq,
      Comparando: "Berlin",
    },
  )

  zachodniopomorskie := poland.Filter(
    dataframe.F{
      Colname:    "County",
      Comparator: series.Eq,
      Comparando: "zachodniopomorskie",
    },
  )

  szczecin := zachodniopomorskie.Filter(
    dataframe.F{
      Colname:    "City",
      Comparator: series.Eq,
      Comparando: "Szczecin",
    },
  )

  fmt.Println(poland)
  fmt.Println(germany)
  fmt.Println(berlin)
  fmt.Println(zachodniopomorskie)
  fmt.Println(szczecin)
}