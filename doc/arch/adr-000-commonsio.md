# commonsio in golang for common file system operation

## context
For common io, file ops, and housekeeping jobs many programming languages have libraries besides their std lib which help reduce verbosity and cognitive load for devs.

Example : java has one from ASF, which identifies few common io and file utilities as follows:
1. FileFilters : filtering files in a directory based on different criteria such as file age, file size, name matching etc.
2. IOUtils : reducing verbosity of the code that reads, writes, and copies files

Here, we will attempt to create similar library for go.

## Decision
While go has excellent std lib support for most of these operations, idea is to find and implement those which are still verbose or missing.

In the first pass of this attempt, focus is on FileFilters utility.

## Status
Accepted.

## Consequences
Implementation in progress.
