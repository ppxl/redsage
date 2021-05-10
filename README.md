# RedSage

Maintain sanity while combining Redmine activity times and Sage project times. This is the deal of RedSage: Sum up correct data into false output data for sanity's sake.

The output times are totally off in favor of easier data input in Sage, that is: every new time slot starts at a full hour. That may mean for a lot of project switches and one official break in a day, Sage-wise will report you to _officially_ work from 08:00 till 22:00 o'clock while in reality you worked from 07:45 till 16:45.

1. Run Redmine query
1. Save .CSV
1. Run RedSage
1. Enter formatted console output in Sage

```bash
# tell RedSage to have a lunch break of 60 minutes (default)
./redsage redmine.csv

# tell RedSage to have a lunch break of 45 minutes 
./redsage -b 45 redmine.csv

# show help topics
./redsage --help
./redsage run --help
```

German CSV quickstart:

Calling RedSage like this reads a Redmine time report .CSV in german locale and joins all pipelines into a single one. Currently, lunch break default to 12:00 o'clock with a duration of 60 minutes.

```
redsage run -s "Gesamtzeit" -c ";" -d "," -i /path/to/timelog-1.csv

Pipeline A    7,50    6,00    4,50    4,50    22,50   
Pipeline B    1,50                            1,50    
ACME                          0,75    0,75    
Pipeline A-joined
2021-05-03      08:00 - 12:00   13:00 - 18:00   
2021-05-04      08:00 - 12:00   13:00 - 15:00   
2021-05-05      08:00 - 12:00   13:00 - 13:30   
2021-05-06      08:00 - 12:00   13:00 - 14:15   
```

Used .CSV:
```csv
Anforderungspipeline;2021-05-03;2021-05-04;2021-05-05;2021-05-06;Gesamtzeit
Pipeline A;7,50;6,00;4,50;4,50;22,50
Pipeline B;1,50;"";"";"";1,50
ACME;"";"";"";0,75;0,75
Gesamtzeit;9,00;6,00;4,50;5,25;24,75
```

## License

MIT