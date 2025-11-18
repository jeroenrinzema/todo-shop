window.onload = function () {
  //<editor-fold desc="Changeable Configuration Block">

  const params = new URLSearchParams(window.location.search);
  window.ui = SwaggerUIBundle({
    url: "/openapi.yaml",
    dom_id: "#swagger-ui",
    deepLinking: true,
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
    plugins: [SwaggerUIBundle.plugins.DownloadUrl],
    layout: "BaseLayout",
  });

  window.ui.initOAuth({
    clientId: params.get("client_id"),
    usePkceWithAuthorizationCodeGrant: true,
  });

  //</editor-fold>
};
