const startPage = function(ctx, config) {
  ctx.get(`${config.baseURL}/`)
}

export default startPage
