# README #

This README will show you how to run the application locally or running it in the Amazon EC2.

## Running locally ##

### Install the code in your home ###

 * cd ~/
 * git clone git@bitbucket.org:csponchiado/rocket.git
 * cd rocket

### Installing the configuration of Solr ###

1. After cloning the project, enter in the directory **solr-config**
2. Run this command: `mvn install`. 
3. It'll create the following structure:
       - `/opt/solr/data[1-4]`              // indexes of the solr simulating each nodes in the cloud
       - `/opt/solr/rocket[1-4]`            // configuration of the solr for each node

### Running Zookeeper ###

1. Install zookeeper
2. Enter in the directory of zookeeper: `cd zookeeper-3.4.6`
3. Run zookeeper : `./bin/zkServer.sh start`

### Running SolrCloud ###

1. Install Solr 4.10.2
2. Enter in the diretory of Solr: `cd solr-4.10.2/example`  
3. Run these commands, each one in one terminal :
	-  `java -DzkHost=localhost:2181 -DnumShards=4 -Dbootstrap_confdir=/opt/solr/rocket1/wikipedia/conf -Dcollection.configName=rocket  -jar -Dsolr.solr.home=/opt/solr/rocket1 -Djetty.port=8080   start.jar`
	-  `java -DzkHost=localhost:2181 -jar -Dsolr.solr.home=/opt/solr/rocket2 -Djetty.port=8081  start.jar`
	-  `java -DzkHost=localhost:2181 -jar -Dsolr.solr.home=/opt/solr/rocket3 -Djetty.port=8082  start.jar`
	-  `java -DzkHost=localhost:2181 -jar -Dsolr.solr.home=/opt/solr/rocket4 -Djetty.port=8083  start.jar`
4. You can check that your solr is running in cloud using this URL: 
	-  `http://localhost:8080/solr/#/~cloud`
5. If you restart your solr, for the port 8080 use this command because the zookeeper has been already configurated
	-  `java -DzkHost=localhost:2181 -jar -Dsolr.solr.home=/opt/solr/rocket1 -Djetty.port=8080  start.jar`

### Running the script of indexing ###

1. Unzip the XML of Wikipedia (http://download.wikimedia.org/enwiki/latest/enwiki-latest-pages-articles.xml.bz2) in this directory: 
	- `~/rocket`
2. Set the Go Path: 
	- `export GOPATH=~/rocket/wikipedia-search/`
3. Run this command to run the indexer: 
	- `./wikipedia-search/bin/indexer -infile=enwiki-latest-pages-articles.xml -solrPort=8080 -solrHost=127.0.0.1 -skip=0`

### Running the script of searching ###

1. Set the Go Path: 
	- `export GOPATH=~/rocket/wikipedia-search/`
2. Run this command:
	- `./wikipedia-search/bin/searcher -solrHost=127.0.0.1 -solrPort=8080   -qf "Title^300 ContributorName^50 ContributorIP^10 Text Title_ngram^50 ContributorName_ngram^3 ContributorIP_ngram^3 Text_ngram" -query=<query>`
		- query: prefix string from ContributorName or Text from wikipedia
		- solrHost: Host of the solr
		- solrPort: Port of the Solr
		- qt: query fields with the boost for each one

## Run using Solr in the Amazon EC2 ##

1. To run the indexer using the solr of Amazon, you have only to change the parameter `-solrHost=<HOST>` from the indexer and searcher
2. The IP of the Solr to use is: 54.94.144.33
3. Let me know when you finish to test to turn off the machine there.

## Questions ##

### a) How did you approach extracting tokens from the article bodies?###

1. First of all, I used the code `xml.NewDecoder` to read the XML as a stream. It's like the SAX API.
2. The Title and ContributorName wasn't necessary to change anything, because their formats was OK to return in the
   response of the searcher. The tokenization and manipulation to create the index was in the Solr.
3. Because the content of the tag Text has a lot of markup from Wikepedia, I had to manipulate it. The text pass in some steps to remove some markups of the text: 
	- Removing references and Citing sources: 
		- `<ref name="name for reference">Use a closing tag</ref>`
		- `{{cite book}} `
	- Removing Links
		- `[[Wikipedia:Manual of Style#Italics]]`
	- Removing others kind of Markup
		- `<math>...</math>`  
		- `<div...>...</div>`
		- `<br>`
		- `#REDIRECT`
	-  Because of the requeriments was to *"Create another tool that takes a **single term prefix** for either a contributor or word that resolves to an article title"*, I removed all the duplicity of the words of the Text Tag so that the index isn't big. If it were necessary to know the position between the terms, it would be necessary indexing the text of Text tag as it is. Reducing the size of the field Text in Solr, it will use less memory and it'll have more performance. 

4. In the solr, I used two types of fields in the schema (`solr-config/src/main/config/solr/wikipedia/conf/schema.xml`). This kind of field manipulate the tokens so that occurs more match in the search. 
	- text_general (general tokenizer)
		- Remove tags HTMLs
		- Remove latin accents
		- Tokenize the words removing all character no alphanumeric
		- Break the words in various formats: Camel, number and words together
		- Lowercase
		- Trim
		- Stopwords (without content)
		- Synonyms (without content)
		- Stemmer
		- Others
	- text_general_ngram (generate the ngram from Title, ContributorName and Text tags)
		- Remove tags HTMLs
		- Remove latin accents
		- Tokenize removing all character not alphanumeric
		- Generate the NGran of the tokens

5. Removing the markup in the text is not perfect. There is some markup like HTML entities that could be removed.
6. The package that manipulate the wiki Text is the wiki.go. I created some unit testing to check if it's correct. 


### b) What would you need to change in your tools to return ranked results?  ###

1. In the Solr, I have used the Extended dismax query parser that is used to process advanced user queries in the solr. I have used the query parameter 'qf' (query field) that has the list of fields and the boost associated to them. You have to pass the fields that you want to search and the weight (importance) of each one.

`go run searcher.go -query=casa -qf "Title^300 ContributorName^50 ContributorIP^10 Text Title_ngram^50 ContributorName_ngram^3 ContributorIP_ngram^3 Text_ngram"`

In this example, if the term is found in the title, it will be better ranked because of the boost of 300 (Title^300). If I use the prefix of a word that exists in the title, it will be boost by 50 (Title_ngram^50).

2. It's possible to cahge the application to use the query parameter mm (minimum should match). It mean the minimum number terms that must match (mm). If I search using 3 terms, all articles that has these terms will be better ranked than the articles that have less terms.


### c) How would your solution change if it needed to be run within 1GB of memory?  If it runs off of disk, describe what/how many disk operations are needed to perform a query. ###

If I have only 1 GB of memory, it'll be necessary divide the index between many nodes. With the use of SolrCloud, I tryed simulate it creating 4 shards so that the index was distributed between these nodes. Each node will require less memory. It's necessary create more shards as the index grows. In my solution, I have used only one zookeeper because the sie of the machine that I got in the amazon, but it's necessary 3 to have better availability. I didn't configurate the replica of the solr, because of the same reason of the zookeeper.

I stored only the title, contributorName and contributorIP inside the index. Because of this requirements to use only 1GB and to show only the title in the results, I would change the schema to put the contributorName and ContributorIP only to be indexed, and not stored, to reduce the use of memory.

In this site has a lot of tips to reduce heap requirements: https://wiki.apache.org/solr/SolrPerformanceProblems

If the user ask for 10 articles for each search, it's necessary 10 disk operation to read the articles in the disc. In the inverted index, each terms point to the ids of the documents that they appears and the frequency that they appear for each document. The documents isn't in the memory of the inverted index, so the solr need to read it in the disc. The solr has 3 kinds of caches: queryCache, filterCache and documentCache. It it's configurated, the 10 articles are read from the disk and put in the cache. For the next search for the same query, the result will be fast.

### d) What parts of the indexing/lookup did you implement yourself? ###

In the indexer, I have prepared the data of the articles before indexing. To be more efficient, I used the goroutine to send 100 articles to the solr at a time. 

In the searcher component, I used a third-party component called Go-Solr because it abstract the connection and objects that is sent and returned from the search. I have only to use the correct parameters to return the articles.

The indexer needs some refactoring to put the manipulation of the text inside of a goroutines to do the work in paralell. It's necessary to do some treatments to handle the errors that can occurs in the connection with the solr.


### e) If you received a diff of changes to articles (additions and removals) what would be needed to merge them into your solution? - Could you do this while still querying results, if so, how? ###

In the solr there are two types for indexing documents: full and partial update. In the full update, it's necessary to send a complete article with all the fields. It's how this solution works. In the partial update, it's possible to send parts of the article, like the Text, ContributorName, ContributorIP, etc. To do the partial update, it's necessary to put the fields that will be updated with the value stored=true in the solr. It's because the solr gets the document stored internally, update the fields and index it again.

It's possible index the documents while querying results. The lucene uses the IndexReader to get a snapshot of the index each time that it's opened. After it, any update in the index, to remove or update it, will not be visible to the index until the indexReader is opened again. In each commit of the index, the indexReader is re-opened and it's updated.  


