# 📰 gross

terminal RSS reader written in Go

## Roadmap

1. Implement checking for new items:
	* Keep old XML in cache
	* If new entries, show them
	* implement a .Read field?
2. Better pager
	* Maybe parse some HTML tags?
	* UGCPolicy + some Unmarshaling may do it
3. Extensiability
	* Implement open-in-browser key
	* Implement way to run custom commands on the link, a command prompt
	* Implement way to load local XML files or use a filter like in newsboat
4. Polishing
	* add viper and standard config dirs
	* add error message if no feeds are supplied
