'use client';

import { useEffect } from 'react';

export function OtelInit() {
  useEffect(() => {
    if (process.env.NODE_ENV === 'development') return;

    // OTel Web SDK initialization
    // Em produção, carregar @opentelemetry/sdk-trace-web e instrumentar fetch
    // Exemplo comentado para evitar bundle excessivo no bootstrap:
    /*
    import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
    import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
    import { FetchInstrumentation } from '@opentelemetry/instrumentation-fetch';
    import { registerInstrumentations } from '@opentelemetry/instrumentation';

    const provider = new WebTracerProvider({
      resource: new Resource({
        [SEMRESATTRS_SERVICE_NAME]: '2pshop-web',
        [SEMRESATTRS_SERVICE_VERSION]: '1.0.0',
        [SEMRESATTRS_DEPLOYMENT_ENVIRONMENT]: process.env.NODE_ENV,
      }),
    });

    provider.addSpanProcessor(new BatchSpanProcessor(new OTLPTraceExporter({
      url: process.env.NEXT_PUBLIC_OTEL_URL,
    })));
    provider.register();

    registerInstrumentations({
      instrumentations: [new FetchInstrumentation()],
    });
    */
  }, []);

  return null;
}
