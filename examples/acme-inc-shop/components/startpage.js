const startPage = function(ctx, config) {
  ctx.get(`${config.baseURL}/`, {
    tag: "start",
  });
}

export default startPage;
