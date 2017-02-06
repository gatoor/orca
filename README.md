# orca trainer

This is the orca trainer repo. To run the trainer you need to configure some base settings:

>> ./trainer --configroot /orca/configuration
             --port 5001
             --cloudprovider aws
             --awsaccesskeyid <access_key>
             --awsaccesskeysecret <access_key_secret>
             --awsregion <region>
             --awsbaseami <awsbaseami>
             --awssshkey <awssshkey>
             --awssshkeypath <keypath>
             --awssgroup <awssgroup>
             --instanceusername ubuntu
             --uri http://localhost:5001


