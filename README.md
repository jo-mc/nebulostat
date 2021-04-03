# nebulostat

Rolling statistics calculator. 

Implementing GNU Scientific Library "Running Statistics" in GO. https://www.gnu.org/software/gsl/doc/html/rstat.html

Requires GO language, and then *go build* (or *go install*) in nebulostat directory.
 
Alternately awk version is available for linux command line usage. (no quantile output as yet)

** Processes one number per line. **

Running statistics are estimates with results converging to truth as number of elements increases. See GNU reference for more information.

Usage: Pipe numbers into program (Linux only)

```
  awk '{ print $3 }' datafile.dat | nebulostat
```
or use with a file argument (Linux or Windows)
```
  nebulostat datafile.dat
```
or awk version
```
  awk -f nebulostat.awk datafile.dat,   
```
or pipe in:
```
 awk '{ print $3 }' datafile.dat | awk -f nebulostat.awk
```

Ouptut:

* The sample Mean
* Std Dev
* The estimated variance
* Largest Value
* Smallest Value
* Estimated Median
* Standard Deviation of the mean.
* Skew
* Kurtosis
* Estimated lower quantile
* Estimated middle quantile (Median)
* Estimated upper quantile
* Number of items
