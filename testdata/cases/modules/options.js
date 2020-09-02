definition.addTarget("testapp.loadtest.party")

definition.setArrivalPhases([{
    duration: 10,
    rate: 42,
    max_clients: MAX_CLIENTS,
  },
]);
