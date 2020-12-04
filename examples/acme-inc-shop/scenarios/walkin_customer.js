import startPage from "../components/startpage.js"

const walkinCustomer = function(d, config) {
  definition.session("walkin customer", function(ctx) {
    startPage(ctx, config);
  })
}

const foo = {
  name: "walkin customer",
  setup: walkinCustomer
}

export default foo
