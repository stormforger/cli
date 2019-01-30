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

definition.session("simple", function(session) {
  session.times(10, function(context) {
    context.get("/" + session.getVar("client_id"));
  });
});
