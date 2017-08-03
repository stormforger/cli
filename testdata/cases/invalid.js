definition.setTarget();

definition.setArrivalPhases([
  { duration: 5 * 60, rate: 50, max_users: 2000 },
  { duration: 15 * 60, rate: 60, max_users: 5000 },
]);

definition.session("invalid", function(session) {
  session.get("/");
});
