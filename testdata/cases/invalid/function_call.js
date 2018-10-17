definition.setTarget("testapp.loadtest.party");

definition.setArrivalPhases([{
  duration: 1 * 60,
  rate: 1,
}]);

definition.session("hello-world", function(session) {
  session.get("/");
  helper(session);
});


function helper(context) {
  context.ge("/helper");
}
