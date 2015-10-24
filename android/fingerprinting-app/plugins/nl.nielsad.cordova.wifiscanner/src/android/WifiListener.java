/**
 * Author: Niels A.D.
 * Project: cordova-wifiscanner (https://github.com/nielsAD/cordova-wifiscanner)
 * License: Apache License v2.0 (http://www.apache.org/licenses/LICENSE-2.0)
 *
 * Java interface to Android's WifiManager.startScan()
 * Based on the Device-Motion Cordova plugin
 */
package nl.nielsad.cordova.wifiscanner;

import java.util.List;

import org.apache.cordova.CordovaWebView;
import org.apache.cordova.CallbackContext;
import org.apache.cordova.CordovaInterface;
import org.apache.cordova.CordovaPlugin;
import org.apache.cordova.PluginResult;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.net.wifi.ScanResult;
import android.net.wifi.WifiManager;

import android.os.Handler;
import android.os.Looper;

public class WifiListener extends CordovaPlugin {

	public enum Status {
		STOPPED,
		STARTING,
		RUNNING,
		ERROR_FAILED_TO_START
	}

	private Status status;
	private CallbackContext callbackContext;
	private WifiManager wifiManager;

	private BroadcastReceiver wifiReceiver = new BroadcastReceiver(){
		public void onReceive(Context c, Intent intent) {
			if (WifiListener.this.status != Status.STOPPED) {
				 WifiListener.this.win();
				 WifiListener.this.wifiManager.startScan();
			}
		}
	};

	private Handler timeoutHandler = new Handler(Looper.getMainLooper());

	private Runnable timeoutRunnable = new Runnable() {
		public void run() {
			WifiListener.this.timeout();
		}
	};

	public WifiListener() {
		this.callbackContext = null;
		this.setStatus(Status.STOPPED);
	}

	private void setStatus(Status status) {
		this.status = status;
	}

	@Override
	public void initialize(CordovaInterface cordova, CordovaWebView webView) {
		super.initialize(cordova, webView);
		this.wifiManager = (WifiManager) cordova.getActivity().getSystemService(Context.WIFI_SERVICE);
	}

	@Override
	public boolean execute(String action, JSONArray args, CallbackContext callbackContext) {
		if (action.equals("get")) {
			this.win();
			return true;
		} else if (action.equals("start")) {
			this.callbackContext = callbackContext;
			if (this.status != Status.RUNNING) {
				this.start();
			}
		}
		else if (action.equals("stop")) {
			if (this.status == Status.RUNNING) {
				this.stop();
			}
		} else {
			return false;
		}

		this.sendPluginResult(new PluginResult(PluginResult.Status.NO_RESULT, ""));
		return true;
	}

	@Override
	public void onDestroy() {
		this.stop();
	}
	
	@Override
	public void onReset() {
		if (this.status == Status.RUNNING) {
			this.stop();
		}
	}

	private PluginResult sendPluginResult(PluginResult res) {
		res.setKeepCallback(true);
		this.callbackContext.sendPluginResult(res);
		return res;
	}
	
	private void startTimeout(long delay) {
		stopTimeout();
		timeoutHandler.postDelayed(timeoutRunnable, delay);
	}
	
	private void stopTimeout() {
		timeoutHandler.removeCallbacks(timeoutRunnable);
	}

	private Status start() {
		if ((this.status == Status.RUNNING) || (this.status == Status.STARTING)) {
			return this.status;
		}

		this.setStatus(Status.STARTING);

		if (wifiManager.isWifiEnabled()) {
			this.cordova.getActivity().registerReceiver(this.wifiReceiver, new IntentFilter(WifiManager.SCAN_RESULTS_AVAILABLE_ACTION));
			this.wifiManager.startScan();
		} else {
			this.setStatus(Status.ERROR_FAILED_TO_START);
		  this.fail(Status.ERROR_FAILED_TO_START, "Wifi is not enabled");
			return this.status;
		}


		this.startTimeout(30000);
		return this.status;
	}

	private void stop() {
		stopTimeout();
		if (this.status != Status.STOPPED) {
			this.cordova.getActivity().unregisterReceiver(wifiReceiver);
			this.setStatus(Status.STOPPED);
		}
	}

	private void timeout() {
		if (this.status == Status.STARTING) {
			this.setStatus(Status.ERROR_FAILED_TO_START);
			this.fail(Status.ERROR_FAILED_TO_START, "WifiManager.startScan() timeout");
		}
	}

	private void fail(Status code, String message) {
		if (callbackContext == null) {
			return;
		}

		JSONObject errorObj = new JSONObject();
		try {
			errorObj.put("code", code.ordinal());
			errorObj.put("message", message);
		} catch (JSONException e) {
			e.printStackTrace();
		}

		this.sendPluginResult(new PluginResult(PluginResult.Status.ERROR, errorObj));
	}

	private void win() {
		if (callbackContext != null) {
			this.sendPluginResult(new PluginResult(PluginResult.Status.OK, this.getScanResults()));
		}
	}

	private JSONArray getScanResults() {
		JSONArray r = new JSONArray();
		List<ScanResult> wifi = this.wifiManager.getScanResults();

		if (wifi != null && wifi.size() > 0)
			for(ScanResult s : wifi)
				try {
					JSONObject o = new JSONObject();
					o.put("BSSID",     s.BSSID);
					o.put("SSID",      s.SSID);
					o.put("level",     s.level);
					r.put(o);
				} catch (JSONException e) {
					e.printStackTrace();
				}

		return r;
	}
}
