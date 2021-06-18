package main

import (
  "fmt"
  "os"
  "encoding/csv"
  "io"
  "log"
  "strconv"
  "time"
  context "context"

  dataframe "github.com/rocketlaunchr/dataframe-go"
  civil "cloud.google.com/go/civil"
)

func parseInt64DefaultZero(s string) int64 {
  res, err := strconv.ParseInt(s, 10, 64)
  if err != nil {
    return 0
  } else { return res }
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

func readIncidencePoland() *dataframe.DataFrame {
  f, err := os.Open("./data/data-pl.csv")
  if err != nil {
    log.Fatal(err)
  }

  reader := csv.NewReader(f)
  reader.Comma = ';'
  reader.FieldsPerRecord = -1

  polandDateNum := dataframe.NewSeriesGeneric("DateNum", civil.Date{}, &dataframe.SeriesInit{})
  polandCounty := dataframe.NewSeriesString("County", &dataframe.SeriesInit{})
  polandCity := dataframe.NewSeriesString("City", &dataframe.SeriesInit{})
  polandCases := dataframe.NewSeriesInt64("Cases", &dataframe.SeriesInit{})
  df := dataframe.NewDataFrame(polandDateNum, polandCounty, polandCity, polandCases)

  for {
        record, err := reader.Read()

        if err == io.EOF {
            break
        }

        if err != nil {
            log.Fatal(err)
        }

        df.Append(nil, parseDateNum(record[0]), record[1], record[2], parseInt64DefaultZero(record[3]))
  }

  return df
}

func readIncidenceGermany() *dataframe.DataFrame {
  f, err := os.Open("./data/data-de.csv")
  if err != nil {
    log.Fatal(err)
  }

  reader := csv.NewReader(f)
  reader.FieldsPerRecord = -1

  date := dataframe.NewSeriesString("Date", &dataframe.SeriesInit{})
  county := dataframe.NewSeriesString("County", &dataframe.SeriesInit{})
  city := dataframe.NewSeriesString("City", &dataframe.SeriesInit{})
  cases := dataframe.NewSeriesInt64("Cases", &dataframe.SeriesInit{})
  caseFreshness := dataframe.NewSeriesInt64("CaseFreshness", &dataframe.SeriesInit{})  
  df := dataframe.NewDataFrame(date, county, city, cases, caseFreshness)

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

        df.Append(nil, record[8], record[2], record[3], parseInt64DefaultZero(record[6]), parseInt64DefaultZero(record[11]))
  }

  return df
}

func main() {
  poland := readIncidencePoland()
  germany := readIncidenceGermany()

  filterBerlin := dataframe.FilterDataFrameFn(func(vals map[interface{}]interface{}, row, nRows int) (dataframe.FilterAction, error) {
    if vals["County"] == "Berlin" {
      return dataframe.KEEP, nil
    }
    return dataframe.DROP, nil
  })

  filterZachodniopomorskie := dataframe.FilterDataFrameFn(func(vals map[interface{}]interface{}, row, nRows int) (dataframe.FilterAction, error) {
    if vals["County"] == "zachodniopomorskie" {
      return dataframe.KEEP, nil
    }
    return dataframe.DROP, nil
  })

  filterSzczecin := dataframe.FilterDataFrameFn(func(vals map[interface{}]interface{}, row, nRows int) (dataframe.FilterAction, error) {
    if vals["City"] == "Szczecin" {
      return dataframe.KEEP, nil
    }
    return dataframe.DROP, nil
  })
  
  berlin, _ := dataframe.Filter(context.Background(), germany, filterBerlin)
  zachodniopomorskie, _ := dataframe.Filter(context.Background(), poland, filterZachodniopomorskie)
  szczecin, _ := dataframe.Filter(context.Background(), zachodniopomorskie.(*dataframe.DataFrame), filterSzczecin)

  end := 10
  opts := dataframe.TableOptions{R: &dataframe.Range{End: &end}}
  fmt.Print(berlin.(*dataframe.DataFrame).Table(opts))
  fmt.Print(zachodniopomorskie.(*dataframe.DataFrame).Table(opts))
  fmt.Print(szczecin.(*dataframe.DataFrame).Table(opts))
}