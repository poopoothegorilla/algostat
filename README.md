# Algostat

Simple app to determine whether the MVRV-Z score is beyond a certain threshold.

It gathers data from the algoexploreapi.

## Usage

```bash
$ go run main.go -ust 1577854800 -lp=.10

2021/05/23 13:32:02 High: 0.63 Low: -0.893
2021/05/23 13:32:02 Recent MVRV: -0.83
```

```bash
$ go run main.go -help

  -high_percentile float
    	percentile to use as a high threshold (default 0.9)
  -hp float
    	percentile to use as a high threshold (shorthand) (default 0.9)
  -low_percentile float
    	percentile to use as a low threshold (default 0.1)
  -lp float
    	percentile to use as a low threshold (shorthand) (default 0.1)
  -unix_start_time int
    	unix timestamp to start measurement
  -ust int
    	unix timestamp to start measurement (shorthand)
```
