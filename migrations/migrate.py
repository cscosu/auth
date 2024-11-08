import os
import random
import csv
with open('seed.csv', mode ='r')as file:
    csvFile = csv.reader(file)
    dest = "002-seed.up.csv"
    dest = open(dest,"w")
    writer = csv.writer(dest, delimiter=',')
    header = next(csvFile)
    writer.writerow(header)
    for lines in csvFile:
        print(lines[4] +", " + lines[5])
        lines[4] = 1731034187 - random.randint(3000,94608000)
        if lines[5] != "NULL":
            lines[4] = 1731034187 - random.randint(3000,94608000)
        writer.writerow(lines)