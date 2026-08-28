'use client';

import { useEffect } from 'react';
import { onLCP, onINP, onCLS, onFCP, onTTFB } from 'web-vitals';

function sendToAnalytics(metric: any) {
  // Envia para o backend de observabilidade ou console
  if (process.env.NODE_ENV === 'development') {
    console.log('[WebVital]', metric.name, metric.value);
  }

  // OTel metric via fetch beacon
  if (typeof navigator !== 'undefined' && 'sendBeacon' in navigator) {
    const payload = JSON.stringify({
      name: `rum.${metric.name.toLowerCase()}`,
      value: metric.value,
      rating: metric.rating,
      delta: metric.delta,
      id: metric.id,
      navigationType: metric.navigationType,
    });
    navigator.sendBeacon('/api/v1/telemetry/rum', payload);
  }
}

export function WebVitals() {
  useEffect(() => {
    onLCP(sendToAnalytics);
    onINP(sendToAnalytics);
    onCLS(sendToAnalytics);
    onFCP(sendToAnalytics);
    onTTFB(sendToAnalytics);
  }, []);

  return null;
}
