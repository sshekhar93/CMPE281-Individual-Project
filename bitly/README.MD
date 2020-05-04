Bitly clone Project Journal

-	YouTube Link: https://youtu.be/7nAKLNN_HhI

-	Strated with understanding the the requiremnets of the project.

-	For that i first checked how url shortening application are working.

-	Understandig  the design of these applications and how to make the scalable.

-	once design part was clear i started learning Go language and how to develope  Web applications in Golang


System Design:

	-	This system consistes of four web application written in Golang
		
		1.	Control Panel Server App: 
			-	This app provide the Rest api to create short links of any url
			-	It uses base Base62 encoding to generate unique short keys for each URL.
			-	Once the short URL are generated, it stores them in MySQL databse.
			-	After storing short url in database it also puts that on a queue in RabbitMQ

		2.	Counter Server:
			-	This app solely provide a counter for the Control Panel servers to use for Base62 encoding.
			-	Purpose of this app it maintain a sync between multiple instances of Control Panel Server so that no two short url are same
			-	This is also written in  Golang

		3.	Link Redirect Server:
			-	This Server is responsible for redirecting short urls to actual websites.
			-	Upon  request it looks up for the short url in local cache, if it is not present in local  cache it retrieves it from NoSQL servers.
			-	Before redirecting itputs a  message on queue of RabbitMQ for trend server
			-	 This is also written in golang

		4.	Trend Server:
			-	This server keeps track of the statics of teh uses of short urls
			-	It achieves that by reading the messages from the queue sent by Link Redirect Server. and updatin the count in the NOSql database
			-	This server is alos responsible for adding short url data in NoSql  database on short url creation. It does this by reading the message sent by Control Panel

		5.	RabbitMQ:
			-	RabbitMQ is used for providing communication channel between the servers theough message queue.

		6.	Kong API Gateway: 
			-	Used KOng to provide public interface for the servers and controlling the flow of request.

Challenges Faced:

	-	Learning a new programming language from start was on of the major challenge
	-	Using Nosql application which was designed to work in same node as containers to deploy on seperate instances and making the to communicate for data replication.
	-	Deploying such a large system on AWS without using any sort of automation. eveything needed to be configured manually. Could be avoided by using environment variables in 		app code insteadof hardcoding.

	-	One ofthe major challenges was to configure Kong for dbless mode for  creating multiple api entrypoints. Due  to lack  of detailed documents i had to figure out most of
		the confures on my own by using hit and trial.

	-	Another technical challenge i faced  was that my MYsql databse would stop responding after sometime. It was due to hign number of connection created by my application  with the database. Fixed it by maintaining one persistent connection and resuing it for all the requests. 
