import splunklib.client as client
import json
import jsonbyteify  

c = client.connect(host="localhost", port="8089", username="ahuynh", password="whiting-annoy-clasp")
c.namespace.owner = "nobody"

"""
c.namespace.app = "TA-FitnessTrackers"
#c.namespace.app = "fitness_for_splunk"

passwords = c.storage_passwords

for password in passwords:
    print password.username
    print password.clear_password
    print password.realm
    print
"""
c.namespace.app = "microsofthealth_data"

passwords = c.storage_passwords

for password in passwords:
    print password.username
    print password.clear_password
    print password.realm
    print



c.namespace.app = 'fitness_for_splunk'

"""
collection_name = "fitbit_tokens"

if collection_name in c.kvstore:
    print "Fitbit kvstore exists!"
    #c.kvstore.delete(collection_name)

# Let's create it and then make sure it exists    
#c.kvstore.create(collection_name)

kvstore = c.kvstore[collection_name]
print json.dumps(kvstore.data.query(), indent=2)

collection_name = "google_tokens"

if collection_name in c.kvstore:
    print "Google kvstore exists!"
    #c.kvstore.delete(collection_name)

# Let's create it and then make sure it exists
#c.kvstore.create(collection_name)

kvstore = c.kvstore[collection_name]
print json.dumps(kvstore.data.query(), indent=2)
"""
collection_name = "microsoft_tokens"

if collection_name in c.kvstore:
    print "Microsoft kvstore exists!"
    #c.kvstore.delete(collection_name)

# Let's create it and then make sure it exists
#c.kvstore.create(collection_name)

kvstore = c.kvstore[collection_name]
print json.dumps(kvstore.data.query(), indent=2)


"""
#userid = "FIEJOK2"
#name = "Andy Huynh"

#jsonn = json.dumps({"_key": userid, "name": name, "token": "test"})
#print jsonn + "\n"

#kvstore.data.insert(jsonn)

#kvstore.data.delete_by_id(userid)
print json.dumps(kvstore.data.query(), indent=2)
#print json.dumps(kvstore.data.query_by_id(userid), indent=2)


#kvstores = c.kvstore



#for kvstore in kvstores:
 #   print kvstore.name
"""
