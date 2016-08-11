#!/bin/bash
# Install Fitness For Splunk to Splunk directory

#echo "Copying fitness_for_splunk"
#cp -rf ./fitness_for_splunk /opt/splunk/etc/apps

#echo "Copying TA-FitnessTrackers"
#cp -rf ./TA-FitnessTrackers /opt/splunk/etc/apps

#echo "Copying microsoft_data"
#cp -rf ./microsoft_data /opt/splunk/etc/apps

echo "Copying Python library dependencies to Splunk Python environment"
cp -rf ./fitness_for_splunk/bin/googleapiclient /opt/splunk2/lib/python2.7/site-packages
cp -rf ./fitness_for_splunk/bin/google_api_python_client-1.5.1.dist-info /opt/splunk2/lib/python2.7/site-packages
cp -rf ./fitness_for_splunk/bin/httplib2 /opt/splunk2/lib/python2.7/site-packages
cp -rf ./fitness_for_splunk/bin/oauth2client /opt/splunk2/lib/python2.7/site-packages
cp -rf ./fitness_for_splunk/bin/pyasn1 /opt/splunk2/lib/python2.7/site-packages
cp -rf ./fitness_for_splunk/bin/pyasn1-0.1.9.dist-info /opt/splunk2/lib/python2.7/site-packages
cp -rf ./fitness_for_splunk/bin/pyasn1_modules /opt/splunk2/lib/python2.7/site-packages
cp -rf ./fitness_for_splunk/bin/pyasn1_modules-0.0.8.dist-info /opt/splunk2/lib/python2.7/site-packages
cp -rf ./fitness_for_splunk/bin/rsa /opt/splunk2/lib/python2.7/site-packages
cp -rf ./fitness_for_splunk/bin/rsa-3.4.2.dist-info /opt/splunk2/lib/python2.7/site-packages
cp -rf ./fitness_for_splunk/bin/simplejson /opt/splunk2/lib/python2.7/site-packages
cp -rf ./fitness_for_splunk/bin/simplejson-3.8.2.dist-info /opt/splunk2/lib/python2.7/site-packages
cp -rf ./fitness_for_splunk/bin/six-1.10.0.dist-info /opt/splunk2/lib/python2.7/site-packages
cp -rf ./fitness_for_splunk/bin/splunklib /opt/splunk2/lib/python2.7/site-packages
cp -rf ./fitness_for_splunk/bin/uritemplate /opt/splunk2/lib/python2.7/site-packages
cp -rf ./fitness_for_splunk/bin/six.py /opt/splunk2/lib/python2.7/site-packages
cp -rf ./fitness_for_splunk/bin/jsonbyteify.py /opt/splunk2/lib/python2.7/site-packages

echo "Done"
