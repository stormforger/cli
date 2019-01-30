definition.setTarget("testapp.loadtest.party");

definition.setArrivalPhases([{
  duration: 5 * 60, // 5min in seconds
  rate: 42.0, // clients per second to launch
  max_clients: 1,
}, ]);

definition.setTestOptions({
  cluster: {
    // sizing: "preflight",
    region: "eu-west-1",
  },
  // dumptraffic: true,
});

definition.session("Black Friday 2019", function (session) {
  session.get("/");
});
