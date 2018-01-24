package com.example.tenntenn.showtoast;

import android.app.Activity;
import android.util.Log;
import android.widget.Toast;

public class ToastViewer {
    public static void showToast(final Activity activity) {
        Log.d("ToastViewer", "showToast");
        activity.runOnUiThread(new Runnable() {
            @Override
            public void run() {
                Toast.makeText(activity, "test toast", Toast.LENGTH_SHORT).show();
            }
        });
    }
}