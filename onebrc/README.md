# 1BRC
This is my solution to the 1 Billion Row Challenge,
as Tweeted by [Gunnar Morling](https://x.com/gunnarmorling/status/1741839724933751238) on Jan 1, 2024

### What is the 1BRC
The One Billion Row Challenge (1BRC), as described on its official 
[Github page](https://github.com/gunnarmorling/1brc/blob/main/README.md), 
is a fun exploration of how far modern Java can be pushed
for aggregating one billion rows from a text file.

### The Challenge
Write a program that retrieves temperature measurement values from a text file 
and calculates the min, mean, and max temperature per weather station. 
There's just one caveat: the file has 1,000,000,000 rows! That's more than 10 GB of data!

The program should print out the min, mean, and max values per station, alphabetically ordered:
```
Hamburg;12.0;23.1;34.2
Bulawayo;8.9;22.1;35.2
Palembang;38.8;39.9;41.0
```

### Implementation Journey
Follow me as I detail my steps and thought process in solving this challenge in this README.md file.
For each commit prefixed with 1brc: in this repo, I will highlight my implementation details for that specific commit,
or specifically what I had done to improve the program in the current commit relative to the last/previous commits.
See <a href="implementation-details.md">Implementation Details</a>

##### Notable links
- [1BRC Offical Website](https://1brc.dev/) 
- [Gunnar Morling's Blogpost](https://www.morling.dev/blog/one-billion-row-challenge/)

### Input Dataset
The input dataset containing 1 billion rows will be seeded from data contained in a 805KB csv file pulled from the
[Official Repo](https://github.com/gunnarmorling/1brc/blob/main/data/weather_stations.csv)
The text file has a simple structure with one measurement value per row:

```
Hamburg;12.0
Bulawayo;8.9
Palembang;38.8
Hamburg;34.2
St. John's;15.2
Cracow;12.6
... etc. ...
```