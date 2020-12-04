import startPage from "../components/startpage.js"
import productPage from "../components/productpage.js"

const orderCustomer = function(d, config) {
  definition.session("order customer", function(ctx) {
    startPage(ctx, config);

    productPage(ctx, config, 4711); // TODO: add datasource to example?
  })
}

const exportData = {
  name: "order customer",
  setup: orderCustomer
};
export default exportData;
