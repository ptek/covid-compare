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

func parseDateGer(s string) civil.Date {
  year,err := strconv.Atoi(s[:4])
  if err != nil { panic(err) }
  month,err := strconv.Atoi(s[5:7])
  if err != nil { panic(err) }
  day,err := strconv.Atoi(s[8:10])
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

  date := series.New([]string{}, series.String, "Date")
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

        normalizedDate := parseDateNum(record[0]).AddDays(-1).String()

        date.Append(normalizedDate)
        county.Append(record[1])
        city.Append(record[2])
        cases.Append(parseIntDefaultZero(record[3]))
  }

  return dataframe.New(date, county, city, cases)
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

        normalizedDate := parseDateGer(record[8]).String()

        date.Append(normalizedDate)
        county.Append(record[2])
        city.Append(record[3])
        cases.Append(parseIntDefaultZero(record[6]))
        caseFreshness.Append(parseIntDefaultZero(record[11]))
  }

  return dataframe.New(date, county, city, cases, caseFreshness)
}

func writeIncidences(df dataframe.DataFrame) {
  f, err := os.Create("./data/data-incidences.csv")
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  err = df.WriteCSV(f)
  if err != nil {
    log.Fatal(err)
  }

  f.Sync()
}

func main() {
  germany := readIncidenceGermany().Filter(
    dataframe.F{Colname: "CaseFreshness", Comparator: series.GreaterEq, Comparando: 0},
  )

  berlin := germany.
    Filter(
      dataframe.F{Colname: "County", Comparator: series.Eq, Comparando: "Berlin"},
    ).
    GroupBy("Date").
    Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM}, []string{"Cases"}).
    Arrange(dataframe.RevSort("Date")).
    Rename("berlin", "Cases_SUM")    

  germany = germany.
    GroupBy("Date").
    Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM}, []string{"Cases"}).
    Arrange(dataframe.RevSort("Date")).
    Rename("germany", "Cases_SUM")

  poland := readIncidencePoland()

  zachodniopomorskie := poland.Filter(
    dataframe.F{Colname: "County", Comparator: series.Eq, Comparando: "zachodniopomorskie"},
  )

  poland = poland.
    Filter(
      dataframe.F{Colname: "County", Comparator: series.Eq, Comparando: "Cały kraj"},
    ).
    GroupBy("Date").
    Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM}, []string{"Cases"}).
    Arrange(dataframe.RevSort("Date")).
    Rename("poland", "Cases_SUM")    

  szczecin := zachodniopomorskie.
    Filter(
      dataframe.F{Colname: "City", Comparator: series.Eq, Comparando: "Szczecin"},
    ).
    GroupBy("Date").
    Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM}, []string{"Cases"}).
    Arrange(dataframe.RevSort("Date")).
    Rename("szczecin", "Cases_SUM")    

  zachodniopomorskie = zachodniopomorskie.
    GroupBy("Date", "County").
    Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM}, []string{"Cases"}).
    Arrange(dataframe.RevSort("Date")).
    Drop("County").
    Rename("zachodniopomorskie", "Cases_SUM")

  fmt.Println(germany)
  fmt.Println(berlin)
  fmt.Println(poland)
  fmt.Println(zachodniopomorskie)
  fmt.Println(szczecin)

  res := germany.InnerJoin(berlin, "Date").InnerJoin(poland, "Date").InnerJoin(zachodniopomorskie, "Date").InnerJoin(szczecin, "Date")
  writeIncidences(res)
}