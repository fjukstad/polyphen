# polyphen
Simple command line tool to submit, check status and download results from
Polyphen2 batch queries. Reads both `vcf` files and already formatted queries. 

# Example 
Submit a batch query (`example.txt` from the [PolyPhen-2 User
Guide](http://genetics.bwh.harvard.edu/pph2/dokuwiki/_media/hg0720.pdf)): 
```
$ polyphen -email=myemail@myisp.com -f example/example.txt -formatted
Polyphen batch query submission completed successfully. You'll get an e-mail at
myemail@myisp.com when the query is completed. Until then you can check the
progress with session ID 78608c2eb3a6d2c3699ab364d9d4205a07b20f2a
```

You can also submit point mutations output from tools such as
[MuTect](http://archive.broadinstitute.org/cancer/cga/mutect): 

```
$ polyphen -email=myemail@myisp.com -f mutect.output.vcf
Polyphen batch query submission completed successfully. You'll get an e-mail at
myemail@myisp.com when the query is completed. Until then you can check the
progress with session ID 78 ... 2a
```


Check the status of the query:
```
$ polyphen -status 78608c2eb3a6d2c3699ab364d9d4205a07b20f2a
Batch query status:
Started Wed Dec  7 15:19:33 2016
Completed Wed Dec  7 15:22:05 2016
```

Download results:

```
$ polyphen -download 78608c2eb3a6d2c3699ab364d9d4205a07b20f2a
Output saved in output
$ ls output/
pph2-full.txt  pph2-log.txt  pph2-short.txt  pph2-snps.txt
```

# Install
- Install [go](http://golang.org) and set it up accordingly (`$GOPATH` etc.) 
- `go get github.com/fjukstad/polyphen`

