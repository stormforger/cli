definition.setTarget();

definition.setArrivalPhases([{
  duration: 1 * 60,
  rate: 1,
}]);

definition.session("invalid", function(session) {
  session.get("/");
});
