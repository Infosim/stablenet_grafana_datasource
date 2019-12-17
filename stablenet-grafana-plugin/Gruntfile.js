/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
module.exports = function(grunt) {

  require('load-grunt-tasks')(grunt);

  grunt.loadNpmTasks('grunt-execute');
  grunt.loadNpmTasks('grunt-contrib-clean');

  grunt.initConfig({

    clean: ["dist"],

    copy: {
      src_to_dist: {
        cwd: 'src',
        expand: true,
        src: ['**/*', '!**/*.ts', '!**/*.scss', "!README_DEVELOPERS.md"],
        dest: 'dist'
      },
    },

    watch: {
      rebuild_all: {
        files: ['src/**/*'],
        tasks: ['default'],
        options: {spawn: false}
      }
    },

    babel: {
      options: {
        sourceMap: true,
        presets:  ['env'],
        plugins: ['transform-object-rest-spread']
      },
      dist: {
        files: [{
          cwd: 'src',
          expand: true,
          src: ['**/*.ts'],
          dest: 'dist',
          ext:'.js'
        }]
      }
      
    }

  });

  grunt.registerTask('default', ['clean', 'copy:src_to_dist', 'babel']);
};
