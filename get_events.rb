#!/usr/bin/env ruby

require 'net/http'
require 'json'
require 'pg'
require 'dotenv'

Dotenv.load

def unescape(str)
  str.gsub(/'/, {"'" => "''"})
end


uri = URI.parse(ENV['DATABASE_URL'])
postgres = PG.connect(uri.hostname, uri.port, nil, nil, uri.path[1..-1], uri.user, uri.password)

def gen_url(path) 
	full_url = URI('https://www.eventbriteapi.com/v3/' + path + 'token=GJ2IG7JJRDXTOF77BOTH')
end

def correct_time(time)
	time[10] = " "
	time
end

time = Time.now.strftime("%Y-%m-%dT%H:%M:%S")

endTime = Time.now + (2*24*60*60)
endTime = endTime.strftime("%Y-%m-%dT%H:%M:%S")

count = 0

url = "events/search/?location.address=london&start_date.range_start=" + time + "&start_date.range_end=" + endTime + "&"

res = Net::HTTP.get(gen_url(url))
jsonObj = JSON.parse(res);
if jsonObj["status_code"] == 429
	puts "Sleeping!"
	sleep(3600)
	res = Net::HTTP.get(gen_url(url))
	jsonObj = JSON.parse(res);
end
number = jsonObj["pagination"]["page_size"]
numberPages = jsonObj["pagination"]["page_count"]
currentPage = 1

loop do 
	for i in 0..number do
		event =  jsonObj["events"][i]
		if event != nil 
			name = unescape(event["name"]["text"])
			description = unescape(event["description"]["text"])
			if description != nil
				description += "\n" + event["url"]
			else 
				description = event["url"]
			end
			venueID = event["venue_id"] 
			venueRes = JSON.parse(Net::HTTP.get(gen_url("venues/" + venueID + "/?")))
			if venueRes["status_code"] == 429
				puts "Sleeping!"
				sleep(3600)
				venueRes = JSON.parse(Net::HTTP.get(gen_url("venues/" + venueID + "/?")))
			end
			latitude = venueRes["address"]["latitude"]
			longitude = venueRes["address"]["longitude"]
			startTime = correct_time(event["start"]["local"])
			endTime = correct_time(event["end"]["local"])
			count = count + 1
			puts count	
			#puts "name: " + name
			#puts "description: " + description
			#puts "latitude: " + latitude
			#puts "longitude: " + longitude
			#puts "start time: " + startTime
			#puts "end time: " + endTime
	
      result = postgres.exec("SELECT * from notes WHERE title='" + name + "' AND starttime='" + startTime + "'")
      if result.num_tuples == 0 
        values = "'" + name + "', " + "'" +  description + "', " + "'" + startTime + "', " + "'" + endTime + "', " + longitude + ", " + latitude
        postgres.exec("INSERT INTO notes(title, comments, starttime, endtime, longitude, latitude) VALUES(" + values+ ")")
      end
		end
	end
	currentPage = currentPage + 1
	res = Net::HTTP.get(gen_url(url + "page=" + currentPage.to_s + "&"))
	break if currentPage > numberPages
end

puts "Processed " + count.to_s + " events!"
