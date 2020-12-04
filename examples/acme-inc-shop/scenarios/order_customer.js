import startPage from "../components/startpage.js"

const orderCustomer = function(d, config) {
  definition.session("order customer", function(ctx) {
    startPage(ctx, config);
  })
}

const foo = {
  name: "order customer",
  setup: orderCustomer
}

export default foo
