definition.addTarget("testapp.loadtest.party")

definition.setArrivalPhases([{
    duration: 60,
    rate: 42,
    max_clients: MAX_CLIENTS,
  },
]);
