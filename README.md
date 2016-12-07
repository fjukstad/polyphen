# polyphen
Simple CLI to submit, check status and download results from Polyphen2 batch
queries. Reads both `vcf` files and already formatted queries. 

# Example 
Submit a batch query: 
```
$ polyphen -email=myemail@myisp.com -f example.txt -formatted
Polyphen batch query submission completed successfully. You'll get an e-mail at myemail@myisp.com when the query is completed. Until then you can check the progress with session ID 78608c2eb3a6d2c3699ab364d9d4205a07b20f2a
```

Check the status of the query
```
$ polyphen -status -id 78608c2eb3a6d2c3699ab364d9d4205a07b20f2a
Batch query status:
Started Wed Dec  7 15:19:33 2016
Completed Wed Dec  7 15:22:05 2016
```

Download results

```
$ polyphen  -download -id 78608c2eb3a6d2c3699ab364d9d4205a07b20f2a
Output saved in output
$ ls output/
pph2-full.txt  pph2-log.txt  pph2-short.txt  pph2-snps.txt
```


