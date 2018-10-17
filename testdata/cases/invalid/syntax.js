definition.setTarget("testapp.loadtest.party");

definition.setArrivalPhases([{
  duration: 1 * 60,
  rate: 1,
});

definition.session("invalid-syntax", function(session) {
  session.get("/");
});
