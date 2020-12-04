const productPage = function(ctx, config, productId) {
  ctx.get(`${config.baseURL}/product/${productId}`, {
    tag: "product",
  });
}

export default productPage;
