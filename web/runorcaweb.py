import os

os.environ["APP_CONFIG_FILE"] = "/Users/michaellawson/development/go/src/gatoor/orca/web/configuration.cfg"
from orcaweb import app
app.run(port=8080, debug=True)

