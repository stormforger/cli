import startPage from "../components/startpage.js"

const walkinCustomer = function(d, config) {
  definition.session("walkin customer", function(ctx) {
    startPage(ctx, config);
  })
}

const exportData = {
  name: "walkin customer",
  setup: walkinCustomer
};
export default exportData;
