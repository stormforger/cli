const config = {
  baseURL: "https://testapp.loadtest.party/shop",
}

definition.addTarget(config.baseURL)

definition.setArrivalPhases([
	{ duration: 10 * 60, rate: 20, },
	{ duration: 10 * 60, rate: 30, },
	{ duration: 10 * 60, rate: 40, },
]);

definition.setTestOptions({
  cluster: { region: "eu-central-1", sizing: "large", },
});

import walkinCustomer from "./scenarios/walkin_customer.js"
walkinCustomer.setup(definition, config);
// definition.addSessionWeight(walkinCustomer.name, 1)

import orderCustomer from "./scenarios/order_customer.js"
orderCustomer.setup(definition, config);
// definition.addSessionWeight(orderCustomer.name, 1)

let weights = {}
weights[walkinCustomer.name] = 1
definition.setSessionWeights(weights)
