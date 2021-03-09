	#! /usr/bin/awk -f
    
BEGIN {
 
rSMin = "NaN"
rSMax = "NaN"
rSN = 0
firstNum = 1
invalid = 0
}
    
{
if ($1 ~ /^[0-9.0-9]+$/)  { # only do numbers

    x = $1 + 0

    rSMin = rSMin + 0
    rSMax = rSMax + 0
    #printf("%s, %s\n",x,rSMin)

    rSMin = (firstNum==1 || x<rSMin ? x : rSMin)
    rSMax = (firstNum==1|| x>rSMax ? x : rSMax)

    firstNum = 0

    /* update mean and variance */
    n1 = rSN
    rSN = rSN + 1
    delta = x - rSM1
    deltaN = delta / rSN
    deltaN2 = deltaN * deltaN
    term1 = delta * deltaN * n1
    n1 = n1 + 1
    rSM1 = rSM1 + deltaN
    rSM4 = rSM4 + term1*deltaN2*(n1*n1-3*n1+3) + 6*deltaN2*rSM2 - 4*deltaN*rSM3
    rSM3 = rSM3 + term1*deltaN*(n1-2) - 3*deltaN*rSM2
    rSM2 = rSM2 + term1
} else {
    invalid += 1
} # if numeric
}
END {

printf("Overall results:\n")
printf("%d, Number of items\n", rSN)
printf("%.2f, The sample Mean\n", rSM1)
printf("%.2f, Std Dev\n", sqrt(rSM2/(rSN-1.0)))
printf("%.2f, The estimated variance\n", (rSM2 / (rSN - 1)))
printf("%.2f, Largest Value\n", rSMax)
printf("%.2f, Smallest Value\n", rSMin)
# printf("%.2f, Estimated Median\n", rs.RQuantResult(medianRq50))
printf("%.2f, Standard Deviation of the mean.\n", ( sqrt(rSM2/(rSN-1.0)) / sqrt(rSN) ))
if (rSN > 0) {
	printf("%.2f, Skew\n", ( (((rSN-1.0)^1.5)/ rSN ) * rSM3 / (rSM2^1.5)))
} else {
	printf("0, Skew")
}
if (rSN > 0) {
	printf("%.2f, Kurtosis\n", (((rSN-1.0) / rSN) * (rSN-1.0)) * rSM4 / (rSM2*rSM2) - 3.0)
} else {
    printf("0, Kurtosis")
}
printf("%s, Invalid data input (discarded)\n",invalid)
}

