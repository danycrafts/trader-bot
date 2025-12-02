from playwright.sync_api import sync_playwright

def verify_stream_controls():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page(viewport={'width': 1280, 'height': 1024})
        try:
            page.goto("http://localhost:5173")
            page.wait_for_selector("#App")
            
            # Check for input box for symbol
            symbol_input = page.locator("input[name='symbol']")
            if symbol_input.is_visible():
                print("Symbol input is visible.")
            else:
                print("Symbol input is NOT visible.")
                exit(1)
                
            # Check for Start Stream button
            start_btn = page.locator("button", has_text="Start Stream")
            if start_btn.is_visible():
                print("Start Stream button is visible.")
            else:
                print("Start Stream button is NOT visible.")
                exit(1)
            
            page.screenshot(path="/home/jules/verification/verification_phase2.png")
            print("Screenshot taken.")
            
        except Exception as e:
            print(f"Error: {e}")
        finally:
            browser.close()

if __name__ == "__main__":
    verify_stream_controls()
