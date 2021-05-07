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
./redsage -b 45 
```

## License

MIT