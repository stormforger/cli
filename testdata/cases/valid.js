definition.setTarget("testapp.loadtest.party");

definition.setArrivalPhases([
  { duration: 5 * 60, rate: 50, max_users: 2000 },
  { duration: 15 * 60, rate: 60, max_users: 5000 },
]);

definition.session("base", function(session) {
  session.get("/users/configuration", {
    gzip: true,
    tag: "user_configuration",
    headers: {
      "Accept": "application/json",
      "X-DemoApp-Token": session.matchedValue("authenticationToken"),
    },
  });

  session.wait(2, { random: true });
});
