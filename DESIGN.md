# Design
In this document, we describe the utilities as we intend to implement one by one.

## FileFilters
In general, applications that produce files as output data, need some kind of housekeeping utility, to delete the files after some intervals.
This is necessary to have sufficient disk available for newer data generated.
Typically, files are deleted based on their age.
A file which has been there on the disk for certain amount of time may not carry significance anymore and could be deleted.
This use case could be of interest to many application, and hence implementing this utility as part of this library.
There is two step approach to this:
1. Find the files with age above or equal to certain value
2. Delete them

The first part could be a feature in itself, where list of files inside a directory made available to clients for reasons other than deletion.

Also, age based filters is one such use case.
There could be cases where filtering is needed based on other criteria.

### Factors for listing the files
1. Depth within a directory upto which listing is needed.
