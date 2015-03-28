'use strict';

var React = require('react-native');
var {
  AppRegistry,
  StyleSheet,
  Text,
  View,
  ListView,
  Image,
  NavigatorIOS
} = React;

var Nav = React.createClass({
  render: function () {
    return (
      <NavigatorIOS
        style={styles.navigator}
        initialRoute={{
          component: Timeline,
          title: 'Timeline',
      }}/>
    );
  }
});

var LoadingView = React.createClass({
  render: function () {
    return (
      <View style={styles.container}>
        <Text>Now Loading...</Text>
      </View>
    )
  },
});

var Timeline = React.createClass({
  getInitialState: function () {
    return {
      tweets: new ListView.DataSource({rowHasChanged: (r1, r2) => r1 !== r2}),
      loaded: false,
    }
  },
  componentDidMount: function() {
    this.fetch();
  },
  render: function () {
    if (!this.state.loaded) {
      return <LoadingView />
    }
    return (
      <ListView
        dataSource={this.state.tweets}
        renderRow={this.renderTweet} />
    )
  },
  renderTweet: function (data, sectionID, rowID) {
    return (
      <View style={styles.tweetContainer}>
        <Image
          source={{uri: data.user.profile_image_url}}
          style={styles.thumbnail} />
        <View style={styles.textContainer}>
          <Text style={styles.textScreenName}>{data.user.screen_name}</Text>
          <Text style={styles.textTweet}>{data.text}</Text>
        </View>
      </View>
    )
  },
  fetch: function () {
    fetch("http://127.0.0.1:8080")
      .then((response) => response.json())
      .then((data) => {
        this.setState({
          tweets: this.state.tweets.cloneWithRows(data),
          loaded: true,
        });
      })
      .done();
  },
});

var styles = StyleSheet.create({
  navigator: {
    flex: 1
  },
  container: {
    flex: 1,
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#FFFFFF',
  },
  thumbnail: {
    width: 80,
    height: 80,
    margin: 2,
  },
  listView: {
    backgroundColor: '#FFFFFF',
  },
  tweetContainer: {
    flex: 1,
    backgroundColor: '#000000',
    borderBottomWidth: 1,
    borderColor: '#E6E6E6',
    flexDirection: 'row',
  },
  textContainer :{
    flex: 1,
  },
  textScreenName: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#FFFFFF',
  },
  textTweet: {
    color: '#FFFFFF',
  },
});

AppRegistry.registerComponent('ReactTweetViewer', () => Nav);
