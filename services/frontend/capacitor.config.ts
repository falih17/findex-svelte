import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
	appId: 'com.findex.attendance',
	appName: 'Findex Attendance',
	webDir: 'build',
	server: {
		androidScheme: 'https',
		cleartext: true, // Allows HTTP access to backend API during development or on private network
		allowNavigation: ['*']
	},
	plugins: {
		Camera: {
			// Auto prompt for native camera permissions
		},
		Geolocation: {
			// High accuracy GPS location
		}
	},
	android: {
		allowMixedContent: true,
		captureInput: true,
		webContentsDebuggingEnabled: true // Enable Chrome DevTools remote debugging for webview
	}
};

export default config;
