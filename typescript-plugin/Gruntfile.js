module.exports = function (grunt) {
    require('load-grunt-tasks')(grunt);

    var pkgJson = require('./package.json');

    grunt.loadNpmTasks('grunt-contrib-clean');
    grunt.loadNpmTasks('grunt-typescript');
    grunt.loadNpmTasks('grunt-execute');

    grunt.initConfig({
        clean: ['dist'],

        copy: {
            src_to_dist: {
                expand: true,
                cwd: 'src',
                src: ['**/*', '!**/*.ts', '!**/*.scss', "!README_DEVELOPERS.md"],
                dest: 'dist'
            },
        },

        typescript: {
            base: {
                src: ['src/*.ts'],
                dest: 'dist',
                options: {
                    module: 'system',
                    target: 'es5',
                    rootDir: 'src/',
                    declaration: true,
                    emitDecoratorMetadata: true,
                    experimentalDecorators: true,
                    sourceMap: true,
                    noImplicitAny: false,
                }
            }
        },

        watch: {
            rebuild_all: {
                files: ['src/**/*'],
                tasks: ['default'],
                options: {spawn: false}
            }
        }
    });

    grunt.registerTask('default', [
        'clean',
        'copy:src_to_dist',
        'typescript',
    ]);
};
