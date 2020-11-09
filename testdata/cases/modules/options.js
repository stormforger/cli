definition.addTarget("testapp.loadtest.party")

definition.setArrivalPhases([{
    duration: 60,
    rate: 42,
    max_clients: 23,
  },
]);
definition.foo(); // Invalid
