curl -v -X POST \
  "https://generativelanguage.googleapis.com/v1beta/models/gemma-3-1b-it:generateContent?key=AIzaSyBUITpP7CzOAgeRyCfGOihHTmn2h5tOqqQ" \
  -H "Content-Type: application/json" \
  -d '{"contents":[{"parts":[{"text":"test"}]}]}'