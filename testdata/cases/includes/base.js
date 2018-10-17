definition.setTarget("testapp.loadtest.party");

definition.setArrivalPhases([{
    duration: 5 * 60,
    rate: 50,
    max_users: 2000
  },
  {
    duration: 15 * 60,
    rate: 60,
    max_users: 5000
  },
]);


//#include ./session_base.js

//#include /Users/basti/Documents/src/golang/src/github.com/stormforger/cli/testdata/cases/includes/helpers.js

function otherHelpers() {
  // ...
}
